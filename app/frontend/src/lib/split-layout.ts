import { get, writable } from "svelte/store"

export const MIN_RATIO = 0.0
export const MAX_RATIO = 1

export const treeRatio = writable(0.75)
export const isResizing = writable(false)

let contentRegion: HTMLDivElement | undefined

export function setContentRegion(node: HTMLDivElement): void {
  contentRegion = node
}

function clampRatio(value: number): number {
  return Math.min(MAX_RATIO, Math.max(MIN_RATIO, value))
}

function updateRatio(clientY: number): void {
  if (!contentRegion) {
    return
  }

  const rect = contentRegion.getBoundingClientRect()
  if (rect.height === 0) {
    return
  }

  const ratio = (clientY - rect.top) / rect.height
  treeRatio.set(clampRatio(ratio))
}

export function handlePointerDown(event: PointerEvent): void {
  if (event.button !== 0 && event.pointerType === "mouse") {
    return
  }

  isResizing.set(true)
  updateRatio(event.clientY)
  event.preventDefault()
}

export function handleWindowPointerMove(event: PointerEvent): void {
  if (!get(isResizing)) {
    return
  }

  updateRatio(event.clientY)
  event.preventDefault()
}

export function handleWindowPointerUp(event: PointerEvent): void {
  if (!get(isResizing)) {
    return
  }

  isResizing.set(false)
  event.preventDefault()
}
