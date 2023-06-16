import { Button, Input, SidePanel } from '@massalabs/react-ui-kit';
import Intl from '../../i18n/i18n';

// no camelcase SidePanel here because it will create conflicts
function Sidepanel({ ...props }) {
  const { url } = props;
  return (
    <SidePanel customClass="border-l border-c-default">
      <div className="flex h-full w-full items-center justify-center">
        <div
          className="flex flex-col justify-center w-[370px] h-fit p-8
                      bg-primary border-dashed border border-c-default
                      "
        >
          <div className="mas-body text-neutral mb-6">
            {Intl.t('store.sidepanel.sidepanel-banner')}
          </div>
          <div className="mas-body2 text-neutral mb-6">
            {Intl.t('store.sidepanel.sidepanel-description', {
              url: url,
            })}
          </div>
          <div className="bg-secondary p-4">
            <div className="mas-menu-active text-neutral mb-3">
              {Intl.t('store.sidepanel.sidepanel-title')}
            </div>
            <div className="mas-caption text-neutral mb-3">
              {Intl.t('store.sidepanel.sidepanel-subtitle')}
            </div>
            <Input
              placeholder={Intl.t('store.sidepanel.sidepanel-placeholder')}
              // default value removes the default green border
              defaultValue=""
              // we use cutom class to fit design of the page
              customClass="bg-primary mb-3"
            />
            <Button>Install</Button>
          </div>
        </div>
      </div>
    </SidePanel>
  );
}

export default Sidepanel;
