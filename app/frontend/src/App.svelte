<script lang="ts">



  let resultText = $state("")
  let name: string

  function greet(): void {
    resultText = "Hello"
  }
</script>

<main>
  <div >Press the button</div>
  <div class="input-box" id="input">
    <button class="btn" onclick={greet}>Hello</button>
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
