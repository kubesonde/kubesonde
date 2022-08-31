import { FaDownload } from "react-icons/fa";
import "./Home.scss";
import React from "react";
import { useFilePicker } from "use-file-picker";
import { useNavigate } from "react-router-dom";
import { ProbeOutput } from "src/entities/probeOutput";
export const HomeComponent = () => {
  const navigate = useNavigate();
  const [openFileSelector, { filesContent, clear }] = useFilePicker({
    accept: ".json",
  });

  if (filesContent.length) {
    const data = filesContent[0].content;
    const parsedData: ProbeOutput = JSON.parse(data);
    const filenameRaw = filesContent[0].name.replace(".json", "");
    const filename = filenameRaw.charAt(0).toUpperCase() + filenameRaw.slice(1);
    clear();
    navigate(`/graph/${filesContent[0].name}`, {
      state: { data: parsedData, title: filename },
    });
  }

  return (
    <div className="form-container">
      <div id="file-upload-form" className="uploader">
        <div id="file-upload" />
        <label htmlFor="file-upload" id="file-drag">
          <img id="file-image" src="#" alt="Preview" className="hidden" />
          <div id="start">
            <FaDownload />
            <div>Select a probe output file</div>
            <span
              id="file-upload-btn"
              className="btn btn-primary"
              onClick={() => openFileSelector()}
            >
              Select a file
            </span>
          </div>
        </label>
      </div>
    </div>
  );
};
