import React from "react";
import { useFilePicker } from "use-file-picker";
import { ProbeOutput } from "src/entities/probeOutput";
import { buildEdgesFromProbes, buildNodesFromProbes } from "src/utils/probes";
import { GraphBase } from "src/components/graph/graphBase/GraphBase";

export const GraphJSONUploadComponent: React.FC = () => {
  const [openFileSelector, { filesContent }] = useFilePicker({
    accept: ".json",
  });

  const getGraph = () => {
    const data = filesContent[0].content;
    const parsedData: ProbeOutput = JSON.parse(data);
    const nodes = buildNodesFromProbes(parsedData);
    const edges = buildEdgesFromProbes(parsedData);
    return <GraphBase title="JSON example" nodes={nodes} edges={edges} />;
  };
  return (
    <>
      {filesContent.length ? (
        getGraph()
      ) : (
        <button onClick={() => openFileSelector()}>
          Select Kubesonde state file
        </button>
      )}
    </>
  );
};
