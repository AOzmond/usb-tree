<script lang="ts">
  import { Init } from "./lib/state.svelte"
  import Header from "./lib/components/Header.svelte"
  import Tree from "./lib/components/Tree.svelte"
  import Log from "./lib/components/Log.svelte"

  import "./style/index.scss"

  import { onMount } from "svelte"

  onMount(() => {
    Init()
  })
</script>

<div class="layout">
  <Header />
  <Tree />
  <Log />
</div>

<style lang="scss">
  :root {
    background: white;
  }
  .layout {
    display: flex;
    flex-direction: column;
    width: 100vw;
    height: 100vh;
    overflow: hidden;
    margin-right: 8px;
  }
</style>
