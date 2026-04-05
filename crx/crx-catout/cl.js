// claude
(() => {
  console.log("cl: Start")
  setInterval(() => {
    console.log("cl: Check")
    if (window.location.search == '?app=claude-code') {
      window.close()
    }
  }, 1000)
})()
