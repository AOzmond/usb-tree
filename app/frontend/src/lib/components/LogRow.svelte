<script lang="ts">
  import { formatSpeed } from "../utilities"
  import { Plus, Minus, Dot } from '@lucide/svelte'


  let { log } = $props()

  const isAdded = $derived(() => log.State === "added")
  const isRemoved = $derived(() => log.State === "removed")
  const isNormal = $derived(() => log.State === "normal")
</script>

<span class="row layout-row {log.State}">
  <div class="left">
    <div class="logstamp">
      <span
        >{new Date(log.Time).toLocaleTimeString([], {
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
          hour12: false,
        })}
      </span>
      {#if isAdded()}
        <Plus />
      {:else if isRemoved()}
        <Minus />
      {:else if isNormal()}
        <Dot />
      {/if}
    </div>
    <div class="logText">
      {log.Text}
    </div>
  </div>
  <div class="right">
    {formatSpeed(log.Speed)}
  </div>
</span>

<style lang="scss">
  .left {
    display: flex;
    flex-direction: row;
    flex: 1 1 auto;
    min-width: 0;
    align-self: center;
    & > :first-child {
      align-self: center;
      white-space: nowrap;
    }
  }
  .logstamp {
    display: flex;
    align-items: center;
  }

  .right {
    white-space: nowrap;
    align-self: center;
  }
</style>
