import plusSmall from '../assets/pictos/plusSmall.svg';
import PlusSmallWhite from '../assets/pictos/plusSmallWhite.svg';
import { UIStore } from '../store/UIStore';
const ManagePluginCard = () => {
  const plus = UIStore.useState((s) =>
    s.theme == 'light' ? plusSmall : PlusSmallWhite,
  );
  return (
    <div
      onClick={() => {
        window.open('/thyra/plugin-manager');
      }}
      className="justify-end center cursor-pointer"
    >
      <h2 className="button text-font underline">Manage plugins ↗</h2>
    </div>
  );
};

export default ManagePluginCard;
