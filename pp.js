const promise1 = new Promise((resolve, reject) => {
  resolve("Success 1");
});

console.log(promise1);
// promise1
//   .then((value) => {})
//   .then((value) => {
//     console.log(value);
//     return new Promise((resolve, reject) => {
//       resolve("Success 2");
//       // console.log("Unresolved");
//     });
//   })
//   .then((value) => console.log(value));

// console.log("End");
