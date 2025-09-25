import type { Action } from "svelte/action";
import { writable } from "svelte/store";

export type TooltipContent = {
  bus: number | null;
  vendorId: string | null;
  productId: string | null;
};

export type TooltipPlacement = "top" | "bottom";

export type TooltipPosition = {
  x: number;
  y: number;
  placement: TooltipPlacement;
};

export type TooltipState = {
  visible: boolean;
  content: TooltipContent | null;
  position: TooltipPosition | null;
};

const initialState: TooltipState = {
  visible: false,
  content: null,
  position: null,
};

const tooltipState = writable<TooltipState>(initialState);

export const tooltip = { subscribe: tooltipState.subscribe };

let hideTimeout: ReturnType<typeof setTimeout> | null = null;

const clearHideTimeout = () => {
  if (!hideTimeout) {
    return;
  }
  clearTimeout(hideTimeout);
  hideTimeout = null;
};

const resetState = () => {
  tooltipState.set(initialState);
};

// showTooltip displays the tooltip with new content at the requested position
export const showTooltip = (
  content: TooltipContent,
  position: TooltipPosition,
) => {
  clearHideTimeout();
  tooltipState.set({
    visible: true,
    content,
    position,
  });
};

export const hideTooltip = (delay = 0) => {
  clearHideTimeout();
  if (delay <= 0) {
    resetState();
    return;
  }

  hideTimeout = setTimeout(() => {
    resetState();
    hideTimeout = null;
  }, delay);
};

// TooltipActionOptions configure how the tooltipTrigger retrieves content and hides
export interface TooltipActionOptions {
  getContent: () => TooltipContent | null;
  hideDelay?: number;
}

// defaultTooltipPosition calculates the pointer-aligned tooltip coordinates
const HEADER_CLEARANCE = 56;
const TOP_OFFSET = 64;
const BOTTOM_OFFSET = 10;

export const defaultTooltipPosition = (
  node: HTMLElement,
  event?: PointerEvent,
): TooltipPosition => {
  const target = (event?.currentTarget as HTMLElement | null) ?? node;
  const rect = target.getBoundingClientRect();
  const preferredTop = rect.top - TOP_OFFSET;
  const placement: TooltipPlacement =
    preferredTop < HEADER_CLEARANCE ? "bottom" : "top";
  const y = placement === "top" ? preferredTop : rect.bottom + BOTTOM_OFFSET;

  return {
    x: rect.left + 40,
    y,
    placement,
  };
};

// tooltipTrigger wires pointer events to the tooltip lifecycle
export const tooltipTrigger: Action<HTMLElement, TooltipActionOptions> = (
  node,
  options,
) => {
  let currentOptions = options;
  let isActive = false;

  const resolveContent = () => currentOptions?.getContent?.() ?? null;

  const scheduleHide = () => {
    const delay = currentOptions?.hideDelay ?? 0;
    hideTooltip(delay);
    isActive = false;
  };

  const handlePointerEnter = (event: PointerEvent) => {
    const content = resolveContent();
    if (!content) {
      isActive = false;
      return;
    }
    showTooltip(content, defaultTooltipPosition(node, event));
    isActive = true;
  };

  const handlePointerLeave = () => {
    if (!isActive) {
      return;
    }
    scheduleHide();
  };

  node.addEventListener("pointerenter", handlePointerEnter);
  node.addEventListener("pointerleave", handlePointerLeave);

  return {
    update(newOptions) {
      currentOptions = newOptions;
    },
    destroy() {
      node.removeEventListener("pointerenter", handlePointerEnter);
      node.removeEventListener("pointerleave", handlePointerLeave);
    },
  };
};
