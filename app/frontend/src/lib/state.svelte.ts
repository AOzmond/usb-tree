import { writable } from "svelte/store";
import { Log, TreeNode } from "$lib/models";
import { EventsOn } from "$wailsjs/runtime/runtime.js";
import { InitFrontend, Refresh } from "$wailsjs/go/main/App";

export const deviceTree = writable<TreeNode[]>([]);
export const deviceLogs = writable<Log[]>([]);

export type CarbonTheme = "g100" | "white";

export const theme = writable<CarbonTheme>("g100");

export function getNextTheme(current: CarbonTheme): CarbonTheme {
  if (current === "g100") {
    return "white";
  }
  return "g100";
}

export function toggleTheme(): void {
  theme.update((mode) => getNextTheme(mode));
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

theme.subscribe((mode) => {
  if (typeof document === "undefined") {
    return;
  }

  const root = document.documentElement;
  root.dataset.theme = mode;
  root.dataset.carbonTheme = mode;
});
