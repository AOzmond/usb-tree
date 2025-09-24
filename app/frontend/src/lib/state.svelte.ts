import { writable } from "svelte/store";
import { Log, TreeNode } from "./models";
import { EventsOn } from "../../wailsjs/runtime/runtime.js";
import { InitFrontend, Refresh } from "../../wailsjs/go/main/App";

export const deviceTree = writable<TreeNode[]>([]);
export const deviceLogs = writable<Log[]>([]);

export type CarbonTheme = "g100" | "white";

export const CARBON_THEME_SEQUENCE: CarbonTheme[] = ["g100", "white"];
const DARK_THEME: CarbonTheme = "g100";
const DEFAULT_THEME: CarbonTheme = "g100";

function resolvedTheme(): CarbonTheme {
  return DEFAULT_THEME;
}

export const theme = writable<CarbonTheme>(resolvedTheme());

theme.subscribe((mode) => {
  if (typeof document === "undefined") {
    return;
  }

  const root = document.documentElement;
  root.dataset.theme = mode;
  root.dataset.carbonTheme = mode;

  const prefersDark = mode === DARK_THEME;
  root.style.colorScheme = prefersDark ? "dark" : "light";
});

export function getNextTheme(current: CarbonTheme): CarbonTheme {
  const index = CARBON_THEME_SEQUENCE.indexOf(current);
  if (index === -1) {
    return DEFAULT_THEME;
  }
  return CARBON_THEME_SEQUENCE[(index + 1) % CARBON_THEME_SEQUENCE.length];
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
