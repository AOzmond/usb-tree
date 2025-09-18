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

// clearHideTimeout used to stop tooltip disappearing when mouseover
const clearHideTimeout = () => {
  if (hideTimeout !== null) {
    clearTimeout(hideTimeout);
    hideTimeout = null;
  }
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

// moveTooltip repositions the tooltip while it remains visible
export const moveTooltip = (position: TooltipPosition) => {
  tooltipState.update((state) =>
    state.visible ? { ...state, position } : state,
  );
};

// scheduleTooltipHide defers hiding to allow hover transitions
export const scheduleTooltipHide = (delay = 120) => {
  clearHideTimeout();
  hideTimeout = setTimeout(() => {
    hideTimeout = null;
    tooltipState.set(initialState);
  }, delay);
};

// cancelTooltipHide prevents a queued hide from firing
export const cancelTooltipHide = () => {
  clearHideTimeout();
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
    x: rect.left + rect.width / 2,
    y: rect.top,
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
    cancelTooltipHide();
    const content = resolveContent();
    if (!content) {
      isActive = false;
      return;
    }
    showTooltip(content, defaultTooltipPosition(node, event));
    isActive = true;
  };

  const handlePointerMove = (event: PointerEvent) => {
    if (!isActive) {
      return;
    }
    moveTooltip(defaultTooltipPosition(node, event));
  };

  const handlePointerLeave = () => {
    scheduleTooltipHide(currentOptions?.hideDelay ?? 120);
    isActive = false;
  };

  node.addEventListener("pointerenter", handlePointerEnter);
  node.addEventListener("pointermove", handlePointerMove);
  node.addEventListener("pointerleave", handlePointerLeave);

  return {
    update(newOptions) {
      currentOptions = newOptions;
    },
    destroy() {
      node.removeEventListener("pointerenter", handlePointerEnter);
      node.removeEventListener("pointermove", handlePointerMove);
      node.removeEventListener("pointerleave", handlePointerLeave);
    },
  };
};
