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
  const links = screen.getAllByRole("menuLink") as HTMLLinkElement[];
  expect(links.length).toBe(1);
  const names = links.map((link: HTMLLinkElement) => link.textContent);
  expect(names).toEqual(["Load example probe"]);

  /* const graphToggle = screen.getByRole("graphLibToggle")
    expect(graphToggle.children.length).toBe(1)
    const tg = screen.getByRole("switch") as HTMLInputElement
    expect(tg.checked).toBe(false)
    fireEvent.click(tg)
    expect(tg.checked).toBe(true)*/
});
