import { render, screen } from "@testing-library/react";
import "@testing-library/jest-dom";
import React from "react";
import { GraphView } from "./GraphView";
import { CompleteExample } from "src/mock/biggerExample";
import { ProbeOutput } from "src/entities/probeOutput";

// GraphBase pulls in Cytoscape and relies on `import.meta.env`, which is hard to
// render inside jsdom. Mock it so the test focuses on GraphView's own wiring
// (data cleanup, table rendering, prop passing). The mock surfaces the title and
// the number of nodes/edges it received so we can assert GraphView built them.
jest.mock("src/components/graph/graphBase/GraphBase", () => ({
  GraphBase: (props: { title: string; nodes: unknown[]; edges: unknown[] }) => (
    <div>
      <div>{props.title}</div>
      <div data-testid="node-count">{props.nodes.length}</div>
      <div data-testid="edge-count">{props.edges.length}</div>
    </div>
  ),
}));

test("renders the graph title and the stats table switch from raw probe data", () => {
  render(<GraphView data={CompleteExample} title="My cluster" />);

  // Title comes from GraphBase, driven by the title prop.
  expect(screen.getByText("My cluster")).toBeInTheDocument();

  // GraphView built nodes and edges from the probe data and passed them down.
  expect(Number(screen.getByTestId("node-count").textContent)).toBeGreaterThan(
    0
  );
  expect(Number(screen.getByTestId("edge-count").textContent)).toBeGreaterThan(
    0
  );

  // The graph stats switch is always rendered.
  expect(screen.getByText(/graph stats/)).toBeInTheDocument();

  // The open-ports-in-pods table switch is rendered because the fixture
  // contains podNetworkingv2 data (proves netinfo2Table wiring runs).
  expect(screen.getByText(/open ports in pods/)).toBeInTheDocument();

  // Declarative configuration switch is rendered because the fixture
  // contains podConfigurationNetworking data.
  expect(
    screen.getByText(/declarative network configuration/)
  ).toBeInTheDocument();
});

test("applies cleanupProbeOutput internally (drops localhost-only entries and renames services)", () => {
  const raw: ProbeOutput = {
    ...CompleteExample,
    // Add a localhost-only listening port that cleanupNetInfo must strip.
    podNetworkingv2: {
      ...CompleteExample.podNetworkingv2,
      podLocal: [{ ip: "127.0.0.1", port: "5000", protocol: "TCP" }],
    },
  };

  render(<GraphView data={raw} title="Cleanup test" />);

  // Renders without crashing and shows the title.
  expect(screen.getByText("Cleanup test")).toBeInTheDocument();
  // The raw input still contains the localhost entry (GraphView must not
  // mutate its input).
  expect(raw.podNetworkingv2.podLocal[0].ip).toEqual("127.0.0.1");
});
