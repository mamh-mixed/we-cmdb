<template>
  <div id="graph" class="home-page"></div>
</template>

<script>
import * as d3 from "d3-selection";
import * as d3Graphviz from "d3-graphviz";
import { getAllLayers, getAllCITypesByLayerWithAttr } from "@/api/server";
import { setHeaders } from "@/api/base.js";
export default {
  data() {
    return {
      layers: [],
      graph: {},
      g: null,
      selectedStatus: ["notCreated", "created"]
    };
  },
  methods: {
    async initGraph(status = []) {
      let origin;
      let edges = {};
      let levels = {};
      let graph;
      let graphviz;

      const initEvent = () => {
        graph = d3.select("#graph");
        graph
          .on("dblclick.zoom", null)
          .on("wheel.zoom", null)
          .on("mousewheel.zoom", null);

        this.graph.graphviz = graph
          .graphviz()
          .zoom(true)
          .scale(1.2)
          .width(window.innerWidth * 0.9)
          .height(window.innerHeight * 0.9)
          .attributer(function(d) {
            if (d.attributes.class === "edge") {
              let keys = d.key.split("->");
              let from = keys[0].trim();
              let to = keys[1].trim();
              d.attributes.from = from;
              d.attributes.to = to;
            }

            if (d.tag === "text") {
              let key = d.children[0].text;
              d3.select(this).attr("text-key", key);
            }
          });
      };

      let layerResponse = await getAllLayers();
      if (layerResponse.status === "OK") {
        let tempLayer = layerResponse.data
          .filter(i => i.status === "active")
          .map(_ => {
            return { name: _.value, layerId: _.codeId, ..._ };
          });
        this.layers = tempLayer.sort((a, b) => {
          return a.seqNo - b.seqNo;
        });
        let ciResponse = await getAllCITypesByLayerWithAttr(
          this.selectedStatus
        );
        if (ciResponse.status === "OK") {
          this.source = [];
          this.source = ciResponse.data;
          this.source.forEach(_ => {
            _.ciTypes &&
              _.ciTypes.forEach(async i => {
                let imgFileSource =
                  i.imageFileId === 0 || i.imageFileId === undefined
                    ? defaultCiTypePNG.substring(0, defaultCiTypePNG.length - 4)
                    : `/cmdb/ui/v2/files/${i.imageFileId}`;
                this.$set(i, "form", {
                  ...i,
                  imgSource: imgFileSource,
                  imgUploadURL: `/cmdb/ui/v2/ci-types/${i.ciTypeId}/icon`
                });
                i.attributes &&
                  i.attributes.forEach(j => {
                    this.$set(j, "form", {
                      ...j,
                      isRefreshable: j.isRefreshable ? "yes" : "no",
                      isDisplayed: j.isDisplayed ? "yes" : "no",
                      isAccessControlled: j.isAccessControlled ? "yes" : "no",
                      isNullable: j.isNullable ? "yes" : "no",
                      isAuto: j.isAuto ? "yes" : "no",
                      isEditable: j.isEditable ? "yes" : "no",
                      searchSeqNo: j.searchSeqNo || 0
                    });
                  });

                this.renderRightPanels();
              });
          });
          let uploadToken = document.cookie
            .split(";")
            .find(i => i.indexOf("XSRF-TOKEN") !== -1);
          setHeaders({
            "X-XSRF-TOKEN": uploadToken && uploadToken.split("=")[1]
          });
          initEvent();
          let nodesString = this.genDOT(ciResponse.data);
          this.loadImage(nodesString);
          this.graph.graphviz.renderDot(nodesString);
          let svg = d3.select("#graph").select("svg");
          svg.append("g").lower();
        }
      }
    },
    renderRightPanels() {
      if (!this.nodeName) return;
      if (!!this.isLayerSelected) {
        this.currentSelectedLayer = this.layers.find(
          _ => _.name === this.nodeName
        );
        this.updatedLayerNameValue = {
          codeId: this.currentSelectedLayer.layerId,
          code: this.currentSelectedLayer.name
        };
        this.handleLayerSelect(this.nodeName);
      } else {
        this.source.forEach(_ => {
          _.ciTypes &&
            _.ciTypes.forEach(i => {
              if (i.ciTypeId && i.ciTypeId === +this.g.id) {
                this.currentSelectedCI = i;
                this.currentSelectedCIChildren =
                  (i.attributes &&
                    i.attributes.sort((a, b) => {
                      return a.displaySeqNo - b.displaySeqNo;
                    })) ||
                  [];
              }
            });
        });
        this.updatedCINameValue = {
          ciTypeId: this.g.id,
          name: this.currentSelectedCI.name
        };
        this.getAllEnumTypes();
        this.getAllEnumCategories();
      }
    },
    genDOT(data) {
      let nodes = [];
      data.forEach(_ => {
        if (_.ciTypes) nodes = nodes.concat(_.ciTypes);
      });
      var dots = [
        "digraph  {",
        'bgcolor="transparent";',
        'Node [fontname=Arial, shape="ellipse", fixedsize="true", width="1.1", height="1.1", color="transparent" ,fontsize=12];',
        'Edge [fontname=Arial, minlen="1", color="#7f8fa6", fontsize=10];',
        'ranksep = 1.1; nodesep=.7; size = "11,8"; rankdir=TB'
      ];
      let layerTag = `node [];`;

      // generate group
      let tempClusterObjForGraph = {};
      let tempClusterAryForGraph = [];
      this.layers.map((_, index) => {
        if (index !== this.layers.length - 1) {
          layerTag += '"' + _.name + '"' + "->";
        } else {
          layerTag += '"' + _.name + '"';
        }

        tempClusterObjForGraph[index] = [`{ rank=same; "${_.name}";`];
        nodes.length > 0 &&
          nodes.forEach((node, nodeIndex) => {
            if (node.layerId === _.layerId) {
              let fontcolor =
                node.status === "notCreated" ? "#10a34e" : "black";
              tempClusterObjForGraph[index].push(
                `"${node.name}"[id="${
                  node.ciTypeId
                }",fontcolor="${fontcolor}", image="${
                  node.form.imgSource
                }.png", labelloc="b"]`
              );
            }
            if (nodeIndex === nodes.length - 1) {
              tempClusterObjForGraph[index].push("} ");
            }
          });
        if (nodes.length === 0) {
          tempClusterObjForGraph[index].push("} ");
        }
        tempClusterAryForGraph.push(tempClusterObjForGraph[index].join(""));
      });
      dots.push(tempClusterAryForGraph.join(""));
      dots.push("{" + layerTag + "[style=invis]}");

      //generate edges
      nodes.forEach(node => {
        node.attributes &&
          node.attributes.forEach(attr => {
            if (attr.inputType === "ref" || attr.inputType === "multiRef") {
              var target = nodes.find(_ => _.ciTypeId === attr.referenceId);
              if (target) {
                dots.push(this.genEdge(nodes, node, attr));
              }
            }
          });
      });
      dots.push("}");
      return dots.join("");
    },
    genEdge(nodes, from, to) {
      const target = nodes.find(_ => _.ciTypeId === to.referenceId);
      let labels = to.referenceName ? to.referenceName.trim() : "";
      return (
        '"' +
        from.name +
        '"->' +
        '"' +
        target.name.trim() +
        '"[taillabel="' +
        labels +
        '", labeldistance=3];'
      );
    },
    loadImage(nodesString) {
      (nodesString.match(/image=[^,]*(files\/\d*|png)/g) || [])
        .filter((value, index, self) => {
          return self.indexOf(value) === index;
        })
        .map(keyvaluepaire => keyvaluepaire.substr(7))
        .forEach(image => {
          this.graph.graphviz.addImage(image, "48px", "48px");
        });
    }
  },
  mounted() {
    this.initGraph();
  }
};
</script>

<style lang="scss" scoped>
.home-page {
  align-items: center;
  display: flex;
  font-size: 75px;
  font-weight: 600;
  height: calc(100vh - 60px);
  justify-content: center;
}
</style>