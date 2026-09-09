import { useLocation } from "react-router-dom";
import { ProbeOutput } from "src/entities/probeOutput";
import React from "react";
import { GraphView } from "src/components/graph/GraphView";

export const GraphFromLocation = () => {
  const { state } = useLocation();
  // @ts-ignore
  const title = state.title as string;
  // @ts-ignore
  const data = state.data as ProbeOutput;
  return <GraphView data={data} title={title} />;
};
