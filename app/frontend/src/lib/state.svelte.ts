import { writable } from "svelte/store";
import { Log, Device, TreeNode } from "./models";
import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime.js";
import { InitFrontend, Refresh } from "../../wailsjs/go/main/App";

export const deviceTree = writable<TreeNode[]>([]);
export const deviceLogs = writable<Log[]>([]);
export type ThemeMode = "light" | "dark";

const prefersDark = (): boolean =>
  typeof window !== "undefined" &&
  typeof window.matchMedia === "function" &&
  window.matchMedia("(prefers-color-scheme: dark)").matches;

const safeGetTheme = (): string | null => {
  if (typeof window === "undefined") {
    return null;
  }
  try {
    return window.localStorage.getItem("theme");
  } catch (error) {
    console.warn("Unable to access localStorage for theme", error);
    return null;
  }
};

const resolvedTheme = (): ThemeMode => {
  const saved = safeGetTheme();
  if (saved === "light" || saved === "dark") {
    return saved;
  }
  return prefersDark() ? "dark" : "light";
};

export const theme = writable<ThemeMode>(resolvedTheme());

const applyTheme = (mode: ThemeMode) => {
  if (typeof document !== "undefined") {
    document.documentElement.dataset.theme = mode;
  }
  if (typeof window !== "undefined") {
    try {
      window.localStorage.setItem("theme", mode);
    } catch (error) {
      console.warn("Unable to persist theme to localStorage", error);
    }
  }
};

theme.subscribe((mode) => {
  applyTheme(mode);
});

export const toggleTheme = () => {
  theme.update((mode) => (mode === "light" ? "dark" : "light"));
};

export function Init(): void {
  EventsOff("treeUpdated");
  EventsOff("logsUpdated");
  EventsOn("treeUpdated", (tree: TreeNode[]) => {
    deviceTree.set(tree);
    console.log(deviceTree);
  });
  EventsOn("logsUpdated", (logs: Log[]) => {
    deviceLogs.set(logs);
  });
  InitFrontend().then();
}

export function refreshDevices(): void {
  Refresh().then();
}
