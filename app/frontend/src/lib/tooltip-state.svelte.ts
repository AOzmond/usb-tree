import type { Action } from "svelte/action"
import { writable } from "svelte/store"

export type TooltipContent = {
  bus: number | null
  vendorId: string | null
  productId: string | null
}

export type TooltipPlacement = "top" | "bottom"

export type TooltipPosition = {
  x: number
  y: number
  placement: TooltipPlacement
}

export type TooltipState = {
  visible: boolean
  content: TooltipContent | null
  position: TooltipPosition | null
}

const initialState: TooltipState = {
  visible: false,
  content: null,
  position: null,
}

const HEADER_CLEARANCE = 56 // px
const TOP_OFFSET = 64 // px
const BOTTOM_OFFSET = 10 // px
const LEFT_OFFSET = 40 // px

const tooltipState = writable<TooltipState>(initialState)

export const tooltip = { subscribe: tooltipState.subscribe }

function resetState(): void {
  tooltipState.set(initialState)
}

// showTooltip displays the tooltip with new content at the requested position
export function showTooltip(content: TooltipContent, position: TooltipPosition): void {
  tooltipState.set({
    visible: true,
    content,
    position,
  })
}

export function hideTooltip(): void {
  resetState()
}

// TooltipActionOptions configure how the tooltipTrigger retrieves content and hides
export interface TooltipActionOptions {
  getContent: () => TooltipContent
}

// TooltipPosition calculates the pointer-aligned tooltip coordinates
export function TooltipPosition(node: HTMLElement, event?: PointerEvent): TooltipPosition {
  const target = (event?.currentTarget as HTMLElement | null) ?? node
  const rect = target.getBoundingClientRect()
  const preferredTop = rect.top - TOP_OFFSET
  const placement: TooltipPlacement = preferredTop < HEADER_CLEARANCE ? "bottom" : "top"
  const y = placement === "top" ? preferredTop : rect.bottom + BOTTOM_OFFSET

  return {
    x: rect.left + LEFT_OFFSET,
    y,
    placement,
  }
}

// tooltipTrigger wires pointer events to the tooltip lifecycle
export function tooltipTrigger(
  node: HTMLElement,
  options: TooltipActionOptions
): ReturnType<Action<HTMLElement, TooltipActionOptions>> {
  let currentOptions = options
  let isActive = false

  function resolveContent(): TooltipContent {
    return currentOptions.getContent()
  }

  function handlePointerEnter(event: PointerEvent): void {
    const content = resolveContent()
    if (!content) {
      isActive = false
      return
    }
    showTooltip(content, TooltipPosition(node, event))
    isActive = true
  }

  function handlePointerLeave(): void {
    if (!isActive) {
      return
    }
    hideTooltip()
    isActive = false
  }

  node.addEventListener("pointerenter", handlePointerEnter)
  node.addEventListener("pointerleave", handlePointerLeave)

  return {
    update(newOptions) {
      currentOptions = newOptions
    },
    destroy() {
      node.removeEventListener("pointerenter", handlePointerEnter)
      node.removeEventListener("pointerleave", handlePointerLeave)
    },
  }
}
