export const createPromise = <T>() => {
  const noop = (data: any) => undefined;
  let resolvePromise: (value: T) => void = noop;
  let rejectPromise: (e: Error) => void = noop;
  const promise = new Promise((resolve, reject) => {
    resolvePromise = resolve;
    rejectPromise = reject;
  });
  return [promise, resolvePromise, rejectPromise] as [
    Promise<T>,
    (value: T) => void,
    (e: Error) => void,
  ];
};
