export const cytoscapeStylesheet = [
    {
        selector: "node",
        style: {
            // The node IS the k8s resource icon, recolored per node (see
            // toCyNode / coloredK8sIcon) — no background tile.
            "background-opacity": 0,
            "background-image": "data(icon)",
            "background-fit": "contain",
            "background-clip": "none",
            "border-width": 0,
            color: "black",
            "min-width": "50px",
            "min-height": "50px",
            "text-wrap": "wrap",
            "text-max-width": "60px",
            "text-events": "yes",
            "text-background-color": "transparent",
            "text-background-opacity": 0,
            "text-background-padding": "0px",
            "text-margin-x": "0px",
            "text-margin-y": "4px",
            "text-valign": "bottom",
            "text-halign": "center",
            "text-overflow-wrap": "break-word",
        }
    },
    {
        selector: 'node[hidden="true"]',
        css: {
            display: "none"
        }
    },
    {
        selector: "node[label]",
        style: {
            label: "data(label)",
            "font-size": "10",
            color: "#14203a",
            "text-halign": "center",
            "text-valign": "bottom",
            "text-margin-x": "0px",
            "text-margin-y": "6px",
            "text-wrap": "wrap",
            "text-max-width": "140px",
            // The deployment/pod color lives in the label's background chip.
            "text-background-color": "data(bg)",
            "text-background-opacity": 1,
            "text-background-padding": "4px",
            "text-background-shape": "roundrectangle",
            "text-border-color": "rgba(20,32,58,0.25)",
            "text-border-width": 1,
            "text-border-opacity": 1,
            "text-overflow-wrap": "whitespace",
        }
    },
    {
        selector: "edge",
        style: {
            width: 1.5,
            "curve-style": "bezier",
            "target-arrow-shape": "triangle",
            "control-point-step-size": 100,
        }
    },
    {
        selector: "edge[label]",
        style: {
            "target-label": "data(label)",
            "font-size": "6",
            "text-background-color": "white",
            "text-background-opacity": 1,
            "text-background-padding": "1px",
            "text-border-color": "black",
            "text-border-style": "solid",
            "text-border-width": 0.5,
            "text-border-opacity": 1,
            "text-rotation": "autorotate",
            'target-text-offset': 35,
        },
    },
    {
        selector: 'edge[hidden="true"]',
        css: {
            display: "none"
        }
    },
    {
        // Allowed and declared/expected.
        selector: 'edge[status="expected"]',
        style: {
            "line-color": "#1f9d6b",
            "target-arrow-color": "#1f9d6b",
        }
    },
    {
        // Allowed but not declared — unexpected open path.
        selector: 'edge[status="unexpected"]',
        style: {
            "line-color": "#c9821a",
            "target-arrow-color": "#c9821a",
        }
    },
    {
        // Denied / blocked connection.
        selector: 'edge[status="denied"]',
        style: {
            "line-color": "#d64550",
            "target-arrow-color": "#d64550",
            "text-border-color": "#d64550",
            "color": "#d64550",
        }
    }
] as Array<cytoscape.StylesheetCSS>;

export const cytoscapeStylesheetPrintMode = [
    {
        selector: "node",
        style: {
            // Print/export: solid white tile with a dark border so the recolored
            // icon reads clearly on a white page (no transparency).
            shape: "round-rectangle",
            "background-color": "#ffffff",
            "background-opacity": 1,
            "background-image": "data(icon)",
            "background-fit": "contain",
            "background-clip": "none",
            padding: "7px",
            "border-width": 1,
            "border-color": "#14203a",
            color: "black",
            "min-width": "50px",
            "min-height": "50px",
            "text-wrap": "wrap",
            "text-max-width": "60px",
            "text-events": "yes",
            "text-background-color": "transparent",
            "text-background-opacity": 0,
            "text-background-padding": "0px",
            "text-margin-x": "0px",
            "text-margin-y": "4px",
            "text-valign": "bottom",
            "text-halign": "center",
            "text-overflow-wrap": "break-word",
        }
    },
    {
        selector: 'node[hidden="true"]',
        css: {
            display: "none"
        }
    },
    {
        selector: "node[label]",
        style: {
            label: "data(label)",
            "font-size": "10",
            color: "#14203a",
            "text-halign": "center",
            "text-valign": "bottom",
            "text-margin-x": "0px",
            "text-margin-y": "6px",
            "text-wrap": "wrap",
            "text-max-width": "140px",
            // The deployment/pod color lives in the label's background chip.
            "text-background-color": "data(bg)",
            "text-background-opacity": 1,
            "text-background-padding": "4px",
            "text-background-shape": "roundrectangle",
            "text-border-color": "rgba(20,32,58,0.25)",
            "text-border-width": 1,
            "text-border-opacity": 1,
            "text-overflow-wrap": "whitespace",
        }
    },
    {
        selector: "edge",
        style: {
            width: 1.5,
            "curve-style": "bezier",
            "target-arrow-shape": "triangle",
            "control-point-step-size": 100,
        }
    },
    {
        selector: "edge[label]",
        style: {
            "target-label": "data(label)",
            "font-size": "6",
            "text-background-color": "white",
            "text-background-opacity": 1,
            "text-background-padding": "1px",
            "text-border-color": "black",
            "text-border-style": "solid",
            "text-border-width": 0.5,
            "text-border-opacity": 1,
            "text-rotation": "autorotate",
            'target-text-offset': 35,
        },
    },
    {
        selector: 'edge[hidden="true"]',
        css: {
            display: "none"
        }
    },
    {
        // Allowed and declared/expected.
        selector: 'edge[status="expected"]',
        style: {
            "line-color": "#1f9d6b",
            "target-arrow-color": "#1f9d6b",
        }
    },
    {
        // Allowed but not declared — unexpected open path.
        selector: 'edge[status="unexpected"]',
        style: {
            "line-color": "#c9821a",
            "target-arrow-color": "#c9821a",
        }
    },
    {
        // Denied / blocked connection.
        selector: 'edge[status="denied"]',
        style: {
            "line-color": "#d64550",
            "target-arrow-color": "#d64550",
            "text-border-color": "#d64550",
            "color": "#d64550",
        }
    }
] as Array<cytoscape.StylesheetCSS>;

