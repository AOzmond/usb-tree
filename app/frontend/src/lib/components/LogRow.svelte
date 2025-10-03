<script lang="ts">
  import { iconByState, formatSpeed, formatTimestamp } from "$lib/utilities"

  import { Dot } from "@lucide/svelte"

  let { log } = $props()

  const LogIcon = iconByState[log.State as keyof typeof iconByState] ?? Dot
</script>

<code class={`log-row ${log.State}`}>
  <div class="left">
    <div class="time-state-wrapper">
      <span>
        {formatTimestamp(log.Time)}
      </span>
      <LogIcon />
    </div>
    <div class="text">
      {log.Text}
    </div>
  </div>
  <div class="speed">
    {formatSpeed(log.Speed)}
  </div>
</code>

<style lang="scss">
  @use "variables.scss" as *;

  .log-row {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: flex-start;

    &.added {
      color: var(--color-added);
    }

    &.removed {
      color: var(--color-removed);
    }

    &.error {
      color: var(--color-error);
    }

    .left {
      display: flex;
      flex-direction: row;
      align-items: center;
      gap: $spacing-03;
      flex-grow: 1;
      flex-shrink: 1;
      flex-basis: auto;
      align-self: center;

      .time-state-wrapper {
        white-space: nowrap;
        display: flex;
        align-items: center;
        gap: $spacing-03;
      }

      .text {
        flex-grow: 1;
        flex-shrink: 1;
        flex-basis: 0%;
        min-width: 0;
      }
    }

    .speed {
      white-space: nowrap;
      align-self: center;
    }
  }
</style>
