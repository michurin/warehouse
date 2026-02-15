// cloudflareaccess.com, WIP
(() => {
  console.log("warp: Start")
  setInterval(() => {
    console.log("warp: Check")
    if (window.location.pathname == '/cdn-cgi/access/refresh-identity' && window.location.search == '?success=true&on_success=callback') {
      window.location.href = 'about:blank'
      // window.close() // Doesn't work :-( "Scripts may close only the windows that were opened by them."
    }
  }, 10000)
})()
