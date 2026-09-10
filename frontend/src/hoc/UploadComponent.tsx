import { FaDownload } from "react-icons/fa";
import "./Upload.scss";
import React, { useEffect } from "react";
import { useFilePicker } from "use-file-picker";
import { useNavigate } from "react-router-dom";
import { ProbeOutput } from "src/entities/probeOutput";

export const UploadComponent = () => {
  const navigate = useNavigate();
  const {
    openFilePicker: openFileSelector,
    filesContent,
    clear,
  } = useFilePicker({
    accept: ".json",
  });

  useEffect(() => {
    let isMounted = true;

    if (filesContent.length && isMounted) {
      const data = filesContent[0].content;
      const parsedData: ProbeOutput = JSON.parse(data);
      const filenameRaw = filesContent[0].name.replace(".json", "");
      const filename =
        filenameRaw.charAt(0).toUpperCase() + filenameRaw.slice(1);

      clear();

      navigate(`/graph/${filesContent[0].name}`, {
        state: { data: parsedData, title: filename },
      });
    }

    return () => {
      isMounted = false;
    };
  }, [filesContent, clear, navigate]);

  return (
    <div className="upload">
      <button
        type="button"
        className="dropzone"
        onClick={() => openFileSelector()}
      >
        <span className="dropzone__icon">
          <FaDownload aria-hidden="true" />
        </span>
        <span className="dropzone__title">Open a probe run</span>
        <span className="dropzone__hint">
          Choose a probe output file to visualize its network graph
        </span>
        <span className="dropzone__cta">Choose file</span>
        <span className="dropzone__note">
          Accepts <code>.json</code> exported by Kubesonde
        </span>
      </button>
    </div>
  );
};
