import { render, screen } from "@testing-library/react";
import App from "./App";
import { BrowserRouter } from "react-router-dom";

test("App runs", () => {
  //const history = createMemoryHistory()
  render(
    <BrowserRouter>
      <App />
    </BrowserRouter>
  );
  //const sidebar = screen.getByRole("sidebar")
  expect(screen.getByText("Kubesonde Viewer")).not.toBeUndefined();
  // Navigation lives on the menu-item rows (onClick), so assert the labels render.
  expect(screen.getByText("Upload a file")).not.toBeUndefined();
  expect(screen.getByText("Load example probe")).not.toBeUndefined();

  /* const graphToggle = screen.getByRole("graphLibToggle")
    expect(graphToggle.children.length).toBe(1)
    const tg = screen.getByRole("switch") as HTMLInputElement
    expect(tg.checked).toBe(false)
    fireEvent.click(tg)
    expect(tg.checked).toBe(true)*/
});
