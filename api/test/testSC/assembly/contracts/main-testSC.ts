import {
  createSC,
  generateEvent,
  fileToByteArray,
} from '@massalabs/massa-as-sdk';

/**
 * Creates a new smart contract with the testSC.wasm file content.
 *
 * @param _ - not used
 */
export function main(_: StaticArray<u8>): i32 {
  const bytes: StaticArray<u8> = fileToByteArray(
    './build/testSC.wasm',
  );

  const testSC = createSC(bytes);
  // This event allows us to get the address of the newly created smart contract to call it in the tests.
  generateEvent(
    `TestSC is deployed at :${testSC.toString()}`,
  );

  return 0;
}
