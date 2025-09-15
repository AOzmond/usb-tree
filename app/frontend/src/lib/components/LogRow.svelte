<script>
  export let log
</script>

<span class="row {log.State}">
  <div class="left">
    {new Date(log.Time).toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false,
    })}
    {log.State}<!-- TODO should be replaced with symbol -->
    {log.Text}
  </div>
  <div class="right">
    <!-- TODO add speed to logs(?) -->
    0.01 Mbps
  </div>
</span>

<style lang="scss">
  .row {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    color: var(--color-text);
  }
  .added {
    color: var(--color-added);
  }
  .removed {
    color: var(--color-removed);
  }
</style>
