import JSZip from 'jszip';

export function validateWebsiteName(websiteName: string): boolean {
  return /^[a-z0-9_.\-~]{3,50}$/.test(websiteName);
}

export function validateWebsiteDescription(description: string): boolean {
  // eslint-disable-next-line no-control-regex
  return /^[\x00-\x7F\xC2-\xF4][\x80-\xBF]*$/u.test(description);
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
