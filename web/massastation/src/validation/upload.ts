import JSZip from 'jszip';

export function validateWebsiteName(websiteName: string): boolean {
  return /^[a-z0-9_.\-~]{3,50}$/.test(websiteName);
}

export function validateWebsiteDescription(description: string): boolean {
  // eslint-disable-next-line no-control-regex
  return isUTF8(description);
}

function isUTF8(value: string) {
  try {
    // Convert the string to a Uint8Array
    const encoder = new TextEncoder();
    const encodedData = encoder.encode(value);

    // Decode the Uint8Array using UTF-8 decoder
    const decoder = new TextDecoder('utf-8');
    decoder.decode(encodedData);

    // If decoding is successful, the value is UTF-8 encoded
    return true;
  } catch (error) {
    // If an error occurs during decoding, the value is not UTF-8 encoded
    return false;
  }
}

export function validateDescriptionLength(description: string): boolean {
  return description.length <= 280;
}

export function validateFileExtension(fileName: string): boolean {
  return fileName.endsWith('.zip');
}

export async function validateFileContent(file: File): Promise<boolean> {
  try {
    const zip = await JSZip.loadAsync(file);
    console.log(zip.file('index.html'));
    return zip.file('index.html') !== null;
  } catch (error) {
    console.error('Error reading or parsing the zip file: ', error);
    return false;
  }
}
