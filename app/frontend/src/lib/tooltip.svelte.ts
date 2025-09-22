import type { Action } from "svelte/action";
import { writable } from "svelte/store";

export interface TooltipContent {
  bus: number | null;
  vendorId: string | null;
  productId: string | null;
}

export interface TooltipPosition {
  x: number;
  y: number;
}

export interface TooltipState {
  visible: boolean;
  content: TooltipContent | null;
  position: TooltipPosition | null;
}

const initialState: TooltipState = {
  visible: false,
  content: null,
  position: null,
};

const tooltipState = writable<TooltipState>(initialState);

export const tooltip = { subscribe: tooltipState.subscribe };

let hideTimeout: ReturnType<typeof setTimeout> | null = null;

// showTooltip displays the tooltip with new content at the requested position
export const showTooltip = (
  content: TooltipContent,
  position: TooltipPosition,
) => {
  tooltipState.set({
    visible: true,
    content,
    position,
  });
};

// TooltipActionOptions configure how the tooltipTrigger retrieves content and hides
export interface TooltipActionOptions {
  getContent: () => TooltipContent | null;
  hideDelay?: number;
}

// defaultTooltipPosition calculates the pointer-aligned tooltip coordinates
export const defaultTooltipPosition = (
  node: HTMLElement,
  event?: PointerEvent,
) => {
  const target = (event?.currentTarget as HTMLElement | null) ?? node;
  const rect = target.getBoundingClientRect();
  return {
    x: rect.left,
    y: rect.top - 64,
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

  const handlePointerEnter = (event: PointerEvent) => {
    const content = resolveContent();
    if (!content) {
      isActive = false;
      return;
    }
    showTooltip(content, defaultTooltipPosition(node, event));
    isActive = true;
  };

  node.addEventListener("pointerenter", handlePointerEnter);

  return {
    update(newOptions) {
      currentOptions = newOptions;
    },
    destroy() {
      node.removeEventListener("pointerenter", handlePointerEnter);
    },
  };
};
