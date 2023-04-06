// The entry file of your WebAssembly module.
import { callerHasWriteAccess, generateEvent } from '@massalabs/massa-as-sdk';
import { Args, stringToBytes } from '@massalabs/as-types';

/**
 * This function is meant to be called only one time: when the contract is deployed.
 *
 * @param _binaryArgs - not used: Arguments serialized with Args
 */
export function constructor(_binaryArgs: StaticArray<u8>): StaticArray<u8> {
  // This line is important. It ensures that this function can't be called in the future.
  // If you remove this check, someone could call your constructor function and reset your smart contract.
  if (!callerHasWriteAccess()) {
    return [];
  }
  generateEvent(`TestSC Constructor called`);
  return [];
}

/**
 * @param binaryArgs - Arguments serialized with Args
 * @returns the emitted event serialized in bytes
 */
export function event(binaryArgs: StaticArray<u8>): StaticArray<u8> {
  const argsDeser = new Args(binaryArgs);
  const id = argsDeser.nextString().expect('Message argument is missing');

  const message = `I'm an event! My id is ${id}`;
  generateEvent(message);
  return stringToBytes(message);
}
