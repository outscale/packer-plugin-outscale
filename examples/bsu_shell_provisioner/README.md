# **Outscale Packer Builder with Shell Provisioning**

This repository contains a **Packer template** for creating an **Outscale Machine Image (OMI)**. The build process provisions an Ubuntu-based VM and configures it using **shell scripts**.

---

## **📌 Overview**
This setup performs the following actions:
1. **Launches an Outscale VM** using Packer.
2. **Runs shell scripts** on the instance to install `octl` and `nginx`.
3. **Creates a new Outscale Machine Image (OMI)**.

---

## **Project Structure**
```
/bsu_shell_provisioner
├── ubuntu-oapi.pkr.hcl             # Packer template using shell provisioners
├── variables.auto.pkvars.hcl       # Variables file for Packer
│── scripts/
│   ├── install_octl.sh             # Installs Outscale octl
│── README.md                       # Documentation
```

---

## **:pushpin: Prerequisites**
Ensure the following tools are installed:

**Packer**  
   ```bash
   packer --version
   ```

---

## **:wrench: Configuration**
### **Set Up Variables (`variables.auto.pkvars.hcl`)**
Edit the `variables.auto.pkvars.hcl` file with the correct values:

```hcl
new_omi_name        = "Ubuntu-24.04-Nginx"
ssh_username        = "outscale"
vm_type             = "tinav6.c4r8p2"
region              = "eu-west-2"
osc_source_image_name = "Ubuntu-24.04-*"
# osc_source_image_id  = "ami-b29cea33"
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
Launch an **Outscale VM**.  
Run the **shell script** to install `octl`.  
Install **nginx** via `apt-get`.  
Save the **new OMI**.

The source OMI is selected dynamically from the most recent Outscale-owned image matching `osc_source_image_name`.
You can also pin a specific image by uncommenting `osc_source_image_id` in the vars file and the matching `source_omi` line in the template.

---

## **Provisioning Scripts**
The **shell script (`scripts/install_octl.sh`)** ensures `octl` is installed:

Additionally, **nginx is installed** using the shell provisioner:

```hcl
provisioner "shell" {
  inline = [
    "sudo apt-get -y update",
    "sudo apt-get -y install nginx",
    "sudo systemctl enable nginx",
    "nginx -v"
  ]
}
```
