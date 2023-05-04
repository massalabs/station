const ManagePluginCard = () => {
  return (
    <div
      onClick={() => {
        window.open('/store');
      }}
      className="justify-end center cursor-pointer"
    >
      <h2 className="button text-font underline">Manage plugins ↗</h2>
    </div>
  );
};

export default ManagePluginCard;
