import { writable } from "svelte/store";
import { Log, Device, TreeNode } from "./models";
import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime.js";
import { InitFrontend, Refresh } from "../../wailsjs/go/main/App";

export const deviceTree = writable<TreeNode[]>([]);
export const deviceLogs = writable<Log[]>([]);

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
