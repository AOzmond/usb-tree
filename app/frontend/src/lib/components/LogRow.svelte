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

<code class={`log-row log-row--${log.State}`}>
  <div class="log-row__primary">
    <div class="log-row__timestamp">
      <span
        >{new Date(log.Time).toLocaleTimeString()}
      </span>
      <LogIcon />
    </div>
    <div class="log-row__text">
      {log.Text}
    </div>
  </div>
  <div class="log-row__speed">
    {formatSpeed(log.Speed)}
  </div>
</code>

<style lang="scss">
@use '../../style/variables.scss';

  .log-row {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: flex-start;
  }

  .log-row__primary {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: variables.$spacing-03;
    flex-grow: 1;
    flex-shrink: 1;
    flex-basis: auto;
    min-width: 0;
    align-self: center;

    & > :first-child {
      white-space: nowrap;
    }
  }

  .log-row__timestamp {
    display: flex;
    align-items: center;
    gap: variables.$spacing-03 - variables.$spacing-01;
  }

  .log-row__text {
    flex: 1;
    min-width: 0;
  }

  .log-row__speed {
    white-space: nowrap;
    align-self: center;
  }

  .log-row--added {
    color: var(--color-added);
  }

  .log-row--removed {
    color: var(--color-removed);
  }

  .log-row--error {
    color: var(--color-error);
  }
</style>
