const promise1 = new Promise((resolve, reject) => {
  resolve("Success!");
});

promise1
  .then((value) => {
    throw "123";
    // expected output: "Success!"
  })
  .catch((value) => {
    console.log(value);
    return 12345;
  })
  .then((value) => {
    console.log(value);
  });
