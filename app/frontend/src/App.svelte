<script lang="ts">
  import { Init } from "$lib/state.svelte"
  import Header from "$lib/components/Header.svelte"
  import Tree from "$lib/components/Tree.svelte"
  import Log from "$lib/components/Log.svelte"
  import Tooltip from "$lib/components/Tooltip.svelte"

  import { onMount } from "svelte"

  onMount(() => {
    Init()
  })
</script>

<div class="app">
  <Header />
  <div class="app__content">
    <Tree />
    <Log />
  </div>
  <Tooltip />
</div>

<style lang="scss">
  @import "./style/variables.scss";
  .app {
    height: 100vh;
    color: var(--color-text);
    background: var(--color-tree-bg);
  }
  .app__content {
    display: flex;
    flex-direction: column;
    height: calc(100vh - 64px);
    background: var(--color-tree-bg);
    flex: 1 1 auto;
  }
</style>
