export type Theme = "light" | "dark";

const STORAGE_KEY = "kubesonde-theme";

// Resolve the theme to use on load: an explicit stored choice wins, otherwise
// fall back to the OS preference.
export const getInitialTheme = (): Theme => {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored === "light" || stored === "dark") {
    return stored;
  }
  const prefersDark =
    typeof window.matchMedia === "function" &&
    window.matchMedia("(prefers-color-scheme: dark)").matches;
  return prefersDark ? "dark" : "light";
};

// Apply a theme to the document and remember the choice.
export const applyTheme = (theme: Theme): void => {
  document.documentElement.setAttribute("data-theme", theme);
  localStorage.setItem(STORAGE_KEY, theme);
};
