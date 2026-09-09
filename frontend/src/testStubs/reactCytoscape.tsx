// Test stub for react-cytoscapejs: the real component instantiates Cytoscape
// against a canvas, which jsdom cannot provide. Rendering a placeholder lets
// the surrounding UI (tables, switches) render in tests.
import React from "react";

const CytoscapeComponent: React.FC<Record<string, unknown>> = () => (
    <div data-testid="cytoscape-stub" />
);

export default CytoscapeComponent;
