import React from "react";
import { CompleteExample } from "src/mock/biggerExample";
import { buildEdgesFromProbes } from "src/utils/probes";
import { buildNodesFromProbes } from "src/utils/probes";
import { GraphBase } from "src/components/graph/graphBase/GraphBase";

export const ExampleGraphComponent: React.FC = () => {
  const nodes = buildNodesFromProbes(CompleteExample);
  const edges = buildEdgesFromProbes(CompleteExample);
  return <GraphBase title="Example" nodes={nodes} edges={edges} />;
};
