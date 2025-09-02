<script lang="ts">
  import {Greet} from '../wailsjs/go/main/App.js'


  let resultText: string = ""
  let name: string

  function greet(): void {
    Greet(name).then(result => resultText = result)
  }
</script>

<main>
  <div >Press the button</div>
  <div class="input-box" id="input">
    <button class="btn" on:click={greet}>Hello</button>
    <div class="result" id="result">{resultText}</div>
  </div>
</main>

<style>

  :root {
    background: black;
    color: white;
  }

  main{
    height: 100vh;
    width: 100vw;
    display:flex;
    flex-direction: column;
    align-items: center;
  }

</style>
