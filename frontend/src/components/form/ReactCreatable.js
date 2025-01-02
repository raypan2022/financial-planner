import CreatableSelect from 'react-select/creatable';

const ReactCreatable = ({ mandatory, name, title, value, onChange, onCreateOption, options, errorDiv, errorMsg}) => {
  return (
    <div className="mb-3">
      <label htmlFor={name} className="form-label">
        {title}{mandatory && <span className="text-danger"> *</span>}
      </label>
      <CreatableSelect
        isClearable
        classNamePrefix={"form-select"}
        value={value.value ? value : null}
        onChange={onChange}
        onCreateOption={onCreateOption}
        options={options}
        placeholder="Select or create an option..."
      />
      <div className={errorDiv}>{errorMsg}</div>
    </div>
  );
};

export default ReactCreatable;
