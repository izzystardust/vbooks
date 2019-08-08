var login = (username, password) => {
  var body = JSON.stringify({
    username: username,
    password: password
  });
  fetch("http://localhost:3000/login", {
    method: "POST",
    body: body
  })
    .then(r => r.json())
    .then(console.log);
  // TODO: store in cookie, https://github.com/js-cookie/js-cookie
};

export { login };
