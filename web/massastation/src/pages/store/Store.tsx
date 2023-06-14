import { FiCodepen, FiGlobe, FiHome } from 'react-icons/fi';
import { Navigator, LayoutStation } from '@massalabs/react-ui-kit';
import MyStation from './MyStation';

export function Store() {
  let navigator = (
    <Navigator
      items={[
        {
          icon: <FiHome />,
          isActive: false,
        },
        {
          icon: <FiCodepen />,
          isActive: true,
        },
        {
          icon: <FiGlobe />,
          isActive: false,
        },
      ]}
      onClickNext={() => console.log('Next clicked')}
      onClickBack={() => console.log('Back clicked')}
    />
  );
  return (
    <LayoutStation navigator={navigator}>
      <div className="mas-banner text-neutral mb-10 mt-24">Modules</div>
      <div className="mas-menu-active mb-4 text-neutral">My Station</div>
      <MyStation />
    </LayoutStation>
  );
}
