<script lang="ts">
  import { formatSpeed } from "../utilities"
  import { Plus, Minus, Dot } from "@lucide/svelte"

  let { log } = $props()

  const iconByState = {
    added: Plus,
    removed: Minus,
    normal: Dot,
  } as const

  const LogIcon = iconByState[log.State as keyof typeof iconByState] ?? Dot
</script>

<code class={`log-row ${log.State}`}>
  <div class="primary">
    <div class="timestamp">
      <span>
        {new Date(log.Time).toLocaleTimeString()}
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
  }

  .primary {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: $spacing-03;
    flex-grow: 1;
    flex-shrink: 1;
    flex-basis: auto;
    min-width: 0;
    align-self: center;

    & > :first-child {
      white-space: nowrap;
    }
  }

  .timestamp {
    display: flex;
    align-items: center;
    gap: $spacing-03 - $spacing-01;
  }

  .text {
    flex: 1;
    min-width: 0;
  }

  .speed {
    white-space: nowrap;
    align-self: center;
  }

  .added {
    color: var(--color-added);
  }

  .removed {
    color: var(--color-removed);
  }

  .error {
    color: var(--color-error);
  }
</style>
