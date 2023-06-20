import { SyntheticEvent, useRef, useState } from 'react';
import Intl from '../../../i18n/i18n';

import {
  Button,
  Input,
  SidePanel,
  TextArea,
  DragDrop,
} from '@massalabs/react-ui-kit';
// import { usePut } from '../../../custom/api';
import {
  validateDescriptionLength,
  validateFileContent,
  validateFileExtension,
  validateWebsiteDescription,
  validateWebsiteName,
} from '../../../validation/upload';
import axios from 'axios';
import { parseForm } from '../../../utils/ParseForm';

interface IFormError {
  websiteName?: string;
  description?: string;
  filename?: string;
}

interface IFormObject {
  websiteName: string;
  description: string;
  filename: string;
}

// interface IUploadRequest {
// }

// interface IUploadResponse {
// }

export default function Upload() {
  const form = useRef(null);
  const [error, setError] = useState<IFormError | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [file, setFile] = useState<File | null>(null);
  // const { mutate: mutableUpload } = usePut<IUploadRequest, IUploadResponse>('websiteUploader/prepare');

  async function validate(e: SyntheticEvent): Promise<boolean> {
    const formObject = parseForm<IFormObject>(e);
    const { websiteName, description } = formObject;

    if (websiteName === '') {
      setError({ websiteName: Intl.t('search.errors.no-website-name') });
      return false;
    }

    if (!validateWebsiteName(websiteName)) {
      setError({ websiteName: Intl.t('search.errors.invalid-website-name') });
      return false;
    }

    if (!validateDescriptionLength(description)) {
      setError({ description: Intl.t('search.errors.description-too-long') });
      return false;
    }

    if (!validateWebsiteDescription(description)) {
      setError({ description: Intl.t('search.errors.invalid-description') });
      return false;
    }

    if (!file) {
      setError({ filename: Intl.t('search.errors.no-file') });
      return false;
    }

    if (!validateFileExtension(file?.name)) {
      setError({ filename: Intl.t('search.errors.invalid-file-extension') });
      return false;
    }

    console.log('validating file content');
    console.log(await validateFileContent(file));

    if (!(await validateFileContent(file))) {
      setErrorMsg(Intl.t('search.errors.invalid-file-content'));
      return false;
    }
    setErrorMsg(null);

    return true;
  }

  function getDefaultWallet(): string {
    return 'buildnet';
  }

  function uploadWebsite(e: SyntheticEvent) {
    const formObject = parseForm<IFormObject>(e);
    const { websiteName, description } = formObject;
    const bodyFormData = new FormData();
    bodyFormData.append('url', websiteName);
    bodyFormData.append('description', description); // Add the website description to the form data
    bodyFormData.append('nickname', getDefaultWallet());
    bodyFormData.append('zipfile', file as File); // we force compiler type of `file` because we know it's not null
    axios({
      url: `/websiteUploader/prepare`,
      method: 'put',
      data: bodyFormData,
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  }

  function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();

    console.log('submitting');
    validate(e).then((isValid) => {
      if (isValid) {
        uploadWebsite(e);
      }
    });
  }

  return (
    <SidePanel customClass="border-l border-c-default bg-secondary">
      <div className="pr-4 m-auto">
        <div
          className={`text-f-primary border-2 border-dashed border-neutral bg-primary p-8`}
        >
          <form ref={form} onSubmit={handleSubmit}>
            <p className="mas-body mb-6">{Intl.t('search.sidebar.title')}</p>
            <div className="flex gap-3 mb-6">
              <p className="mas-body2">{Intl.t('search.sidebar.how-to')}</p>
              <h3 className="mas-h3 underline cursor-pointer">
                howtouploadwebsite.com
              </h3>
            </div>
            <div className="bg-secondary rounded-lg p-4 mb-6">
              <p className="mas-menu-active mb-3">
                {Intl.t('search.sidebar.your-website')}
              </p>
              <p className="mas-caption mb-3">
                {Intl.t('search.sidebar.your-website-description')}
              </p>
              <div className="pb-3">
                <Input
                  defaultValue=""
                  name="websiteName"
                  placeholder={Intl.t('search.inputs.website-name')}
                  customClass="mb-3 bg-primary"
                  error={error?.websiteName}
                />
              </div>
              <TextArea
                defaultValue=""
                name="description"
                placeholder={Intl.t('search.inputs.website-description')}
                error={error?.description}
              />
            </div>
            <div className="mb-6">
              <DragDrop
                onFileLoaded={(file) => setFile(file)}
                placeholder="something to inform"
                allowed={['zip']}
              />
            </div>
            <Button type="submit">{Intl.t('search.buttons.upload')}</Button>
          </form>
          {errorMsg && <p className="mas-body pt-4 text-s-error">{errorMsg}</p>}
        </div>
      </div>
    </SidePanel>
  );
}
