const promise1 = new Promise((resolve, reject) => {
  return resolve(123);
});

promise1
  .then((value) => {
    console.log(value);
    return new Promise((resolve, reject) => {
      resolve("Success 2");
    });
  })
  //   .then((value) => {
  //     console.log(value);
  //   })
  //   .then((value) => {
  //     console.log(value);
  //     return new Promise((resolve, reject) => {
  //       resolve("Success 3");
  //     });
  //   })
  //   .then((value) => {
  //     console.log(value);
  //     throw "123";
  //   })
  //   .catch((value) => {
  //     console.log("Fail");
  //   })
  //   .finally(() => {
  //     console.log("Final 1");
  //     return new Promise((resolve, reject) => {
  //       console.log("Sending");
  //       resolve("Success 56");
  //     });
  //   })
  //   .finally(() => {
  //     console.log("Final 2");
  //   })
  .then()
  .then((value) => {
    console.log(value);
  });
