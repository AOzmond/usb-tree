import { writable } from "svelte/store";
import { Log, TreeNode } from "./models";
import { EventsOn } from "../../wailsjs/runtime/runtime.js";
import { InitFrontend, Refresh } from "../../wailsjs/go/main/App";

export const deviceTree = writable<TreeNode[]>([]);
export const deviceLogs = writable<Log[]>([]);
export type ThemeMode = "light" | "dark";

function prefersDark(): boolean {
  return (
    typeof window !== "undefined" &&
    typeof window.matchMedia === "function" &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  );
}

function resolvedTheme(): ThemeMode {
  return prefersDark() ? "dark" : "light";
}

export const theme = writable<ThemeMode>(resolvedTheme());

theme.subscribe((mode) => {
  if (typeof document === "undefined" || typeof window === "undefined") {
    return;
  }
  document.documentElement.dataset.theme = mode;
});

export function toggleTheme(): void {
  theme.update((mode) => (mode === "light" ? "dark" : "light"));
}

export function Init(): void {
  EventsOn("treeUpdated", (tree: TreeNode[]) => {
    deviceTree.set(tree);
  });
  EventsOn("logsUpdated", (logs: Log[]) => {
    deviceLogs.set(logs);
  });
  InitFrontend().then();
}

export function refreshDevices(): void {
  Refresh().then();
}
