<script lang="ts">
  import { formatSpeed } from "../utilities"
  let { log } = $props()
</script>

<span class="row {log.State} ">
  <div class="left">
    <div>
      {new Date(log.Time).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
        hour12: false,
      })}
      {log.State}<!-- TODO should be replaced with symbol -->
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
  .row {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    color: var(--color-text);
    align-items: flex-start;
  }
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
  .right {
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
    color: orange;
  }
</style>
