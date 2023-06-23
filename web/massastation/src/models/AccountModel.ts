type keyPairObject = {
  nonce: string;
  privateKey: string;
  publicKey: string;
  salt: string;
};

export type AccountObject = {
  address: string;
  balance: string;
  candidateBalance: string;
  keyPair: keyPairObject;
  nickname: string;
};
