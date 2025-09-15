<script lang="ts">
  import { Init } from "./lib/state.svelte"
  import Header from "./lib/components/Header.svelte"
  import Tree from "./lib/components/Tree.svelte"
  import Log from "./lib/components/Log.svelte"

  import { onMount } from "svelte"

  onMount(() => {
    Init()
  })
</script>

<div class="layout">
  <Header />
  <div class="content">
    <Tree />
    <Log />
  </div>
</div>

<style lang="scss">
  :global {
    @import "./style/index.scss";
  }
  @import "./style/variables.scss";

  :global(body) {
    background: var(--color-tree-bg);
    color: var(--color-text);
  }

  .layout {
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
    height: 100vh;
    overflow: hidden;
    color: var(--color-text);
    background: var(--color-tree-bg);
  }
  .content {
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    background: var(--color-tree-bg);
    height: 100%;
  }
</style>
