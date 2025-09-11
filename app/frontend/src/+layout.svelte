<script lang="ts">
  import { Init } from "./lib/state.svelte"
  import Header from "./lib/components/Header.svelte"
  import Tree from "./lib/components/Tree.svelte"
  import Log from "./lib/components/Log.svelte"
  Init()
</script>

<Header />
<Tree />
<Log />

<style>
  :root {
    background: white;
    font-size: 16;
  }


</style>
