<script>
  export let log
</script>

<span class={log.State}>
  {log.Time}
  {log.Text}
  {log.State}
</span>

<style lang="scss">
  .added {
    color: green;
  }
  .removed {
    color: red;
  }
</style>
