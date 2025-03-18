# **Outscale Packer Builder with Ansible**

This repository contains a **Packer template** for creating an **Outscale Machine Image (OMI)**. The build process provisions an Ubuntu-based VM and configures it using **Ansible**.

---

## **:pushpin: Overview**
This setup performs the following actions:
1. **Launches an Outscale VM** using Packer.
2. **Runs an Ansible Playbook** on the instance to install Apache.
3. **Creates a new Outscale Machine Image (OMI)**.

---

## **Project Structure**
```
/bsu_ansible_provisioner
│── packer/
│   ├── ubuntu-apache.pkr.hcl             # Packer template using Ansible
│   ├── variables.auto.pkvars.hcl         # Variables file for Packer
│── ansible/
│   ├── playbook.yml                      # Ansible playbook (Installs Apache)
│── logs/
│   ├── packer_build.log                   # Log file created after build
│── README.md                              # Documentation
```

---

## **:pushpin: Prerequisites**
Ensure the following tools are installed:
1. **Packer**  
   ```bash
   packer --version
   ```
2. **Ansible** (Required for local execution)  
   ```bash
   sudo apt update && sudo apt install -y ansible
   ```

---

## **:wrench: Configuration**
### **Set Up Variables (`variables.auto.pkvars.hcl`)**
Edit the `variables.auto.pkvars.hcl` file with the correct values:

```hcl
new_omi_name        = "Ubuntu-24.04-Apache"
ssh_username        = "outscale"
vm_type             = "tinav6.c4r8p2"
region              = "eu-west-2"
osc_source_image_id = "ami-860c2495"
```

---

## **:pushpin: Running the Build**
### **Initialize Packer**
Run the following command inside the **Packer directory** to initialize required plugins:
```bash
packer init .
```

### **Validate the Configuration**
Before building, check if the configuration is valid:
```bash
packer validate .
```

### **Build the Outscale Machine Image**
```bash
packer build .
```
This process will:
- Launch an **Outscale VM**.
- Run the **Ansible Playbook** to install **Apache**.
- Save the **new OMI**.

---

## **Ansible Playbook**
The **Ansible Playbook (`ansible/playbook.yml`)** installs **Apache** on the VM.

---

## **:pushpin: Logging**
After the build, a log file is created in the **logs directory**:
```bash
cat packer_build.log
```
Example output:
```
Packer build completed at Tue Mar 18 14:52:36 UTC 2025
```